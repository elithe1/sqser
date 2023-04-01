# Input Plugins

This section is for developers who want to create new inputs or understand the architchture better. Sqser is entirely
plugin driven.

## Input Plugin Guidelines

- A plugin must conform to the [models.Input][] interface. All four sqser commands has to be implemented by the inputs.
- Input Plugins should call `inputs.Add` in their `init` function to register themselves.
- To be available within Sqser itself, plugins must register themselves using a file
  in `github.com/elithe1/sqser/plugins/inputs/all/all` named according to the plugin name.
- Each plugin requires a file called `sample.yaml` containing the sample configuration for the plugin in yaml format.
  Please consult the [Config][] page for the latest guidelines.
- Each plugin `README.md` file should include the `sample.yaml` file.
- Input plugin should have a parameterised logger with the input details.
- Inputs can be sync or async. As currently async input is the only one sqser has it's gonna return 200 directly to the
  caller if data parsed correctly. Then do all processing and outputs asynchronously in a new go-routine.

[models.Input]: http://github.com/elithe1/sqser/blob/master/models/input.go#L7-L7

[Config]: http://github.com/elithe1/sqser/blob/master/docs/CONFIG.md
