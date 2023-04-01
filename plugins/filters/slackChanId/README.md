# Slack channel id Filter Plugin

When using slack input, we know in which channel the command was invoked. This plugin often allows you might have an
integration for slack with different channels for different environments (listed in environments substring) and then
this plugin will filter in only the relevant queues for the channel The plugin will try to find what is the envirnment
that belongs to the chnnelId the message is sent from. Then it will allow only queues that belong to this environment.
name

## Configuration

Values should have a list of:

- **name**: environment name (matches the environments subsctring) and
- **id**: the id of the slack channel
