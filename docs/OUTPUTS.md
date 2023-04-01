# Output Plugins

This section is for developers who want to create new outputs or understand the architecture better. Sqser is entirely
plugin driven.

## Output Plugin Guidelines

- A plugin must conform to the [models.Output][] interface. All four sqser commands has to be implemented by the
  outputs.
- Output Plugins should call `outputs.Add` in their `init` function to register themselves.
- To be available within Sqser itself, plugins must register themselves using a file
  in `github.com/elithe1/sqser/plugins/outputs/all/all` named according to the plugin name.
- Each plugin requires a file called `sample.yaml` containing the sample configuration for the plugin in yaml format.
  Please consult the [Config][] page for the latest guidelines.
- Each plugin `README.md` file should include the `sample.yaml` file.- Input plugin should have a parameterised logger
  with the input details.

[models.Output]: http://github.com/elithe1/sqser/blob/master/models/output.go#L7-L7

[Config]: http://github.com/elithe1/sqser/blob/master/docs/CONFIG.md
