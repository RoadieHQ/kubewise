package config

type Config struct {
  Namespace string

  // The # prefix should be included in the channel name
  SlackChannel string
  SlackToken string
}
