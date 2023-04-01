# Filter Plugins

This section is for developers who want to create new filters or understand the architecture better. Sqser is entirely
plugin driven.

## Filter Plugin Guidelines

- A plugin must conform to the [models.Filter][] interface.
- Input Plugins should call `inputs.Add` in their `init` function to register themselves.
- To be available within Sqser itself, plugins must register themselves using a file
  in `github.com/elithe1/sqser/plugins/filters/all/all` named according to the plugin name.
- Each plugin requires a file called `sample.conf` containing the sample configuration for the plugin in yaml format.
  Please consult the [Config][] page for the latest style guidelines.
- Each plugin `README.md` file should include the `sample.yaml` file.- Input plugin should have a parameterised logger
  with the input details.
- Input plugin should have a parameterised logger with the input details.
- Inputs can be sync or async. As currently async input is the only one sqser has it's gonna return 200 directly to the
  caller if data parsed correctly. Then do all processing and outputs asynchronously in a new go-routine.

[models.Input]: http://github.com/elithe1/sqser/blob/master/models/filter.go#L7-L7

[Config]: http://github.com/elithe1/sqser/blob/master/docs/CONFIG.md
