![KubeWise Mark and Name](./assets/kubewise-name-and-mark-487x127.png)

![Go workflow status](https://github.com/larderdev/kubewise/workflows/Go/badge.svg)

KubeWise is a notifications bot for Helm. It notifies your team chat whenever a Helm chart is installed,
upgraded or uninstalled in your Kubernetes cluster.

![A demo of KubeWise posting Slack messages as ZooKeeper is installed, upgraded and uninstalled](./assets/kubewise-demo.gif)

# Supported Chat Apps

| Logo | Name | Supported |  |
| ------------- | ------------- | ------------ | ------- |
| ![Slack mark](./assets/slack-mark-50x50.png)  | [Slack](https://slack.com)  | ✅ | [Get started](#slack) |
| ![Google Chat mark](./assets/googlechat-mark-50x50.png)  | [Google Hangouts Chat](https://gsuite.google.com/products/chat/)  | ✅ | [Get started](#google-hangouts-chat) |
| ![Microsoft Teams mark](./assets/ms-teams-mark-50x50.png) | [Microsoft teams](https://products.office.com/en-us/microsoft-teams/group-chat-software) | ⏳ |  |
| ![Flock mark](./assets/flock-mark-50x50.jpg) | [Flock](https://flock.com/) | ⏳ |  |
| ![Mattermost mark](./assets/mattermost-mark-50x50.png) | [Mattermost](https://mattermost.com) | ⏳ |  |
|  | [Twist](https://twist.com) | ⏳ |  |
|  | [Telegram](https://telegram.org) | ⏳ |  |

📣 [Get notified when your chosen chat app is supported.](https://forms.gle/bWJAaaiYArMJ9hrYA)

# Getting Started

In general, the getting started process has two steps:

1. Create a bot in your team chat application.
2. Install KubeWise, passing it an API token for the bot.

Sensitive tokens are stored securely in Kubernetes secrets. No data is ever sent to an external API (other
than your chosen team chat app obviously).

## Slack

### How it looks
![Slack sample](./assets/slack-sample-935x422.png)

### Step 1: Create the bot
 1. Create a [Slack Bot](https://my.slack.com/services/new/bot).
    - username: `kubewise`
    - name: `KubeWise`
    - icon: [Use This](https://raw.githubusercontent.com/larderdev/kubewise/master/assets/kubewise-mark-blue-512x512.png)
 2. Save it and grab the API token.
 3. Invite the Bot into your channel by typing `/invite @kubewise` in your Slack channel.
 4. Install KubeWise in your Kubernetes cluster. See below.

### Step 2: Install KubeWise
```
kubectl create namespace kubewise
helm repo add larder https://charts.larder.dev
helm install kubewise larder/kubewise --namespace kubewise --set handler=slack --set slack.token="<api-token>" --set slack.channel="#<channel>"
```

That's it! From now on, Helm operations will result in a message in your chosen Slack channel.

## Google Hangouts Chat

### How it looks
![Google Hangouts Chat sample](./assets/googlechat-sample-915x605.png)

### Step 1: Create the bot
 1. Open [Hangouts Chat](https://chat.google.com/) in your browser.
 2. Go to the room to which you want to add a bot.
 3. From the dropdown menu at the top of the page, select "Configure webhooks".
 4. Under Incoming Webhooks, click ADD WEBHOOK.
 5. Name the new webhook `KubeWise` and set the the Avatar URL to `https://raw.githubusercontent.com/larderdev/kubewise/master/assets/kubewise-mark-blue-512x512.png`.
 6. Click SAVE.
 7. Copy the URL listed next to your new webhook in the Webhook Url column. You will need this later.
 8. Click outside the dialog box to close.

### Step 2: Install KubeWise
```
kubectl create namespace kubewise
helm repo add larder https://charts.larder.dev
helm install kubewise larder/kubewise --namespace kubewise --set handler=googlechat --set googlechat.webhookUrl="<webhook-url>"
```

# Using KubeWise from outside a cluster

It is easy to use KubeWise from outside your Kubernetes cluster. It will pick up your local
`kubectl` configuration and use it to speak to your cluster.

You will need to compile the go binary from source. For example,

```
# Clone and compile the binary
git clone git@github.com:larderdev/kubewise.git
cd kubewise
go build

# Run it against a cluster
env KW_HANDLER=slack KW_SLACK_CHANNEL="#<channel>" KW_SLACK_TOKEN="<api-token>" kubewise
```
