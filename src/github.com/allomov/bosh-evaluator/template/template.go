package template

func Interpolate() error {
	tpl := boshtpl.NewTemplate(opts.Args.Manifest.Bytes)

	vars := opts.VarFlags.AsVariables()
	op := opts.OpsFlags.AsOp()
	evalOpts := boshtpl.EvaluateOpts{
		ExpectAllKeys:     opts.VarErrors,
		ExpectAllVarsUsed: opts.VarErrorsUnused,
	}

	if opts.Path.IsSet() {
		evalOpts.PostVarSubstitutionOp = patch.FindOp{Path: opts.Path}

		// Printing YAML indented multiline strings (eg SSH key) is not useful
		evalOpts.UnescapedMultiline = true
	}

	bytes, err := tpl.Evaluate(vars, op, evalOpts)
	if err != nil {
		return err
	}

	c.ui.PrintBlock(string(bytes))
}
