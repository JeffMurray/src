The golang json package is great, and its use of reflection to encode and decode structs stands on its own.

Where it needs a little help, I think, is a more freestyle way to get at the data.

In javascript accessig some known variiable nested deep in wrapper objects is trivial:

var key = lib3.apps.darwin.config.access.key;

In golang, if you have the same variable in a map[string]interface{}, trivial does not come to mind.

With package ezjson you can get the key from the map like this:

key, err := ezjson.GetString(config, "lib3", "apps", "darwin", "config", "access", "key")

It has a number of paired Get and Set functions for various golang data types that might get passed in json.

It also helps check for many type conversion problems and has helper functions for encoding and decoding.