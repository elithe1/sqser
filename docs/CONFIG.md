# Configuration

Sqser's configuration file is written using [Yaml][] and is composed of sections: [substrings][], and four sections
describing plugins:
[inputs][], [filters], [enrichers][], and [outputs][].

You can choose any and as many plugins you want and the configuration file is the only source of truth of the plugins
you will have loaded in your Sqser instance startup.

## Substrings

As all queues have some kind of naming convention, we use this section to help sqser to manage the dlqs by setting some
predefined substrings:

- **dlq**: The substring that is marking a queue as a Deadletter queue. Main usages of this substring:
    - Throw errors in case someone is trying to fetch items or move items from a queue that is not marked as a dlq.
    - In order to transfer items from dlq to queue with the move command. CUrrently, the queue name is expected to be
      the same name as dlq just without the dlq substring.

- **environments**: A list of strings which are environments you have resources deployed to. A resolved environment name
  from queue name can be later on used in the enricher plugins. It can help in cases where the enrichers are connected
  to different accounts (eg  [Logzio][] enricher)

## Plugins

Sqser plugins are divided into 4 types: [inputs][], [outputs][],
[filters][], and [enrichers][].

Unlike the `Substrings`, any plugin can be defined multiple times and each instance will run independently. This allows
you to have plugins defined with differing configurations as needed within a single Sqser process.

Each plugin has a unique set of configuration options, reference the sample configuration for details. Each plugin has
the format of a name which is a string and then name of the plugin and values which is an object which holds the
different options the plugin requires (structure defferes between plugins)

### Input Plugins

Sqser is based on a simple http server. Input plugins allow to initially parse the requests into Sqser structs received
by the http server, return success or errors.

#### Examples

```yaml
inputs:
  - name: slack
    async: true

```

### Output Plugins

Output plugins are in charge of outputting the result of the command.

#### Examples

```yaml
outputs:
  - name: slack
```

### Filter Plugins

Filter plugins are used by the list command which has the purpose to list all non-empty dlqs and their number of visible
items.

#### Examples

This filter will keep only non-empty queues that has the "dlq" substring. It will not report back all the non-dlq
non-empty queues

```yaml
filters:
  - name: substringAllowList
    values:
      - dlq
```

### Enricher Plugins

Enricher plugins are used to enrich and data from different sources for the get item command. Often you might want to
construct and add relevant logs data oe some data from your db or anything else to help debugging

#### Examples

If the order processors are applied matters you must set order on all involved processors:

```yaml
enrichers:
  - name: logzio
    values:
      accounts:
        - name: staging
          id: 499492
        - name: production
          id: 499512
      enrichFields:
        timeStamp: timestamp
        searchField: traceId
```

[YAML]: https://yaml.org/

[substrings]: #Substrings

[Logzio]:  http://github.com/elithe1/sqser/blob/master/plugins/enrichers/logzio/logzio.go#L87-L87

[plugins]: #plugins

[inputs]: #input-plugins

[outputs]: #output-plugins

[filters]: #filter-plugins

[enrichers]: #enricher-plugins