package cmd

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type Cmd struct {
	BoshOpts BoshOpts
	Opts     interface{}

	deps boshcmd.BasicDeps
}

func NewCmd(BoshOpts BoshOpts, opts interface{}, deps boshcmd.BasicDeps) Cmd {
	return Cmd{BoshOpts, opts, deps}
}

type cmdConveniencePanic struct {
	Err error
}

func (c Cmd) Execute() (cmdErr error) {
	defer func() {
		if r := recover(); r != nil {
			if cp, ok := r.(cmdConveniencePanic); ok {
				cmdErr = cp.Err
			} else {
				panic(r)
			}
		}
	}()

	c.configureUI()
	c.configureFS()

	deps := c.deps

	switch opts := c.Opts.(type) {
	case *WriteOpts:
		return NewWriteCmd(deps.UI, c.director()).Run()
	case *ReadOpts:
		return NewReadCmd(c.config(), deps.UI).Run()
	default:
		return fmt.Errorf("Unhandled command: %#v", c.Opts)
	}
}

func (c Cmd) configureUI() {
	c.deps.UI.EnableTTY(c.BoshOpts.TTYOpt)

	if !c.BoshOpts.NoColorOpt {
		c.deps.UI.EnableColor()
	}

	if c.BoshOpts.JSONOpt {
		c.deps.UI.EnableJSON()
	}

	if c.BoshOpts.NonInteractiveOpt {
		c.deps.UI.EnableNonInteractive()
	}
}

func (c Cmd) configureFS() {
	tmpDirPath, err := c.deps.FS.ExpandPath("~/.bosh/tmp")
	c.panicIfErr(err)

	err = c.deps.FS.ChangeTempRoot(tmpDirPath)
	c.panicIfErr(err)
}

func (c Cmd) config() cmdconf.Config {
	config, err := cmdconf.NewFSConfigFromPath(c.BoshOpts.ConfigPathOpt, c.deps.FS)
	c.panicIfErr(err)

	return config
}

func (c Cmd) session() Session {
	return NewSessionFromOpts(c.BoshOpts, c.config(), c.deps.UI, true, true, c.deps.FS, c.deps.Logger)
}

func (c Cmd) director() boshdir.Director {
	director, err := c.session().Director()
	c.panicIfErr(err)

	return director
}

func (c Cmd) deployment() boshdir.Deployment {
	deployment, err := c.session().Deployment()
	c.panicIfErr(err)

	return deployment
}

func (c Cmd) directorAndDeployment() (boshdir.Director, boshdir.Deployment) {
	sess := c.session()

	director, err := sess.Director()
	c.panicIfErr(err)

	deployment, err := sess.Deployment()
	c.panicIfErr(err)

	return director, deployment
}

func (c Cmd) releaseProviders(algos []boshcrypto.Algorithm) (boshrel.Provider, boshreldir.Provider) {
	indexReporter := boshui.NewIndexReporter(c.deps.UI)
	blobsReporter := boshui.NewBlobsReporter(c.deps.UI)
	releaseIndexReporter := boshui.NewReleaseIndexReporter(c.deps.UI)

	digestCalculator := c.deps.DigestCalc(algos)

	releaseProvider := boshrel.NewProvider(
		c.deps.CmdRunner, c.deps.Compressor, digestCalculator, c.deps.FS, c.deps.Logger)

	releaseDirProvider := boshreldir.NewProvider(
		indexReporter, releaseIndexReporter, blobsReporter, releaseProvider,
		digestCalculator, c.deps.CmdRunner, c.deps.UUIDGen, c.deps.Time, c.deps.FS, c.deps.Logger)

	return releaseProvider, releaseDirProvider
}

func (c Cmd) releaseManager(director boshdir.Director) ReleaseManager {
	relProv, relDirProv := c.releaseProviders([]boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1})

	releaseDirFactory := func(dir DirOrCWDArg) (boshrel.Reader, boshreldir.ReleaseDir) {
		releaseReader := relDirProv.NewReleaseReader(dir.Path)
		releaseDir := relDirProv.NewFSReleaseDir(dir.Path)
		return releaseReader, releaseDir
	}

	releaseWriter := relProv.NewArchiveWriter()

	createReleaseCmd := NewCreateReleaseCmd(releaseDirFactory, releaseWriter, c.deps.FS, c.deps.UI)

	releaseArchiveFactory := func(path string) boshdir.ReleaseArchive {
		return boshdir.NewFSReleaseArchive(path, c.deps.FS)
	}

	uploadReleaseCmd := NewUploadReleaseCmd(
		releaseDirFactory, releaseWriter, director, releaseArchiveFactory, c.deps.CmdRunner, c.deps.FS, c.deps.UI)

	return NewReleaseManager(createReleaseCmd, uploadReleaseCmd)
}

func (c Cmd) blobsDir(dir DirOrCWDArg) boshreldir.BlobsDir {
	_, relDirProv := c.releaseProviders([]boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1})
	return relDirProv.NewFSBlobsDir(dir.Path)
}

func (c Cmd) releaseDir(dir DirOrCWDArg) boshreldir.ReleaseDir {
	_, relDirProv := c.releaseProviders([]boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1})
	return relDirProv.NewFSReleaseDir(dir.Path)
}

func (c Cmd) panicIfErr(err error) {
	if err != nil {
		panic(cmdConveniencePanic{err})
	}
}
