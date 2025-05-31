📰 Slack News Bot
A simple Slack bot built in Go that fetches the latest tech headlines from Inc42 and responds to user commands in Slack. Also integrates with Make.com via a webhook for automation scenarios.

🚀 Features
Responds to ping with pong

Greets users with hello

Fetches the latest 5 news headlines using RSS feed (news from go)

Sends news headlines to a Make.com webhook (news from make)

Logs command events for debugging

📦 Dependencies
Go 1.18+

mmcdole/gofeed — for RSS parsing

shomali11/slacker — for Slack bot framework

⚙️ Setup Instructions
1. Clone the repository
bash
Copy
Edit
git clone https://github.com/yourusername/slack-news-bot.git
cd slack-news-bot
2. Install dependencies
bash
Copy
Edit
go mod tidy
3. Set environment variables
You'll need to create a Slack bot and get your tokens.

bash
Copy
Edit
export SLACK_BOT_TOKEN='xoxb-your-bot-token'
export SLACK_APP_TOKEN='xapp-your-app-token'
You can also use a .env file or any Go environment loading package if preferred.

4. Run the bot
bash
Copy
Edit
go run main.go
5. Slack Commands
Command	Description
ping	Health check — replies with pong
hello	Replies with a friendly greeting
news from go	Fetches news and displays them in Slack
news from make	Fetches news and sends it to a Make.com webhook

🔗 Make.com Integration
The news from make command sends the fetched news to this webhook:

bash
Copy
Edit
https://hook.eu2.make.com/kev3jirwkvz4attya8y36oav619du9e5
Make sure to configure your Make.com scenario to handle the JSON payload:

json
Copy
Edit
{
  "event": "news_requested",
  "text": "Here are the latest headlines:\n1. ..."
}
🧪 Example Output
bash
Copy
Edit
Here are the latest headlines:
1. *Startup raises $10M in funding*
https://inc42.com/news/startup-raises-10m/

2. *New fintech regulation announced*
https://inc42.com/news/new-fintech-rules/
🛠 To Do
Add unit tests

Allow keyword-based news search

Enhance error handling/logging

Add support for multiple RSS sources

📄 License
MIT License. See LICENSE for details.
