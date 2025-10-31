# yamlark 
### scriping YAML files with starlark

The idea is simple: read YAML config files into your starlark script with `yaml.read()`.

Now it's a dictionary that can be easily manipulated.

Once it's done, convert it back to a YAML string using `yaml.dumps()`.

Additionally, you can:
- read and write files with `file.read()` and `file.write()`.
- read TOML files with `toml.read()`.

See the `example` directory.
