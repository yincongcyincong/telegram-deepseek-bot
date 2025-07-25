# ✨ Telegram & Discord DeepSeek Bot

This project is a cross-platform chatbot powered by the **DeepSeek LLM**, supporting both **Telegram** and **Discord**. It comes with a variety of built-in commands, including image and video generation, balance checking, conversation clearing, and more.

## 🚀 Starting in Discord Mode

You can launch the bot in Discord mode using the following command:

```bash
./telegram-deepseek-bot-darwin-amd64 \
  -discord_bot_token=xxx \
  -deepseek_token=sk-xxx \
  -volc_ak=xxx \
  -volc_sk=xxx \
  -vol_token=xxx
```

### Parameter Descriptions:

* `discord_bot_token`: Your Discord Bot Token (required)
* `deepseek_token`: Your DeepSeek API Token (required)
* `volc_ak` / `volc_sk`: Volcano Engine Access Key and Secret Key (required for `/photo` and `/video` commands)
* `vol_token`: Token for using Volcano Engine video capabilities

other usage see this [doc](https://github.com/yincongcyincong/telegram-deepseek-bot))

## 💬 How to Use

### Private Chat with the Bot

You can directly chat with the bot via private message.    
<img width="400" alt="image" src="https://github.com/user-attachments/assets/6d8ded05-8454-4946-9025-bdd4bb7f8dbb" />


Supported commands:

* `/photo`: Generate an image. ⚠️ Requires `volc_ak` and `volc_sk`
<img width="400" alt="image" src="https://github.com/user-attachments/assets/325a7fab-6cc5-4088-870c-bab3b3c184d8" />
<img width="400" alt="image" src="https://github.com/user-attachments/assets/a117963c-1c21-4f8b-a8f3-ab7ec217040d" />

<img width="400" alt="image" src="https://github.com/user-attachments/assets/ba0eb926-7924-4c58-bc61-7cff522bd71c" />
<img width="400" alt="image" src="https://github.com/user-attachments/assets/bef94980-4498-4eba-b4b5-bd5531816009" />

* `/video`: Generate a video. ⚠️ Requires `volc_ak` and `volc_sk`
<img width="400" alt="image" src="https://github.com/user-attachments/assets/24bdde29-685c-4af7-8834-873dbc14b84f" />
<img width="400" alt="image" src="https://github.com/user-attachments/assets/b9e85a58-58fe-4e45-ab44-52b73bcaea59" />
  
* `/balance`: Check the remaining quota of your DeepSeek Token
<img width="400" alt="image" src="https://github.com/user-attachments/assets/bb20e8fd-470f-4c70-b584-abc1fb5855d2" />

  
* `/state`: View the current chat state (including model info and system prompts)
<img width="400" alt="image" src="https://github.com/user-attachments/assets/bf57f0fa-add1-4cb2-8e82-7bd484a880b8" />
   
* `/clear`: Clear the current conversation context
<img width="400" alt="image" src="https://github.com/user-attachments/assets/ebba556f-267a-4052-a3a3-3eab019eb4f4" />

  

### Group Chat Mode

In a group chat, you can talk to the bot by mentioning it using `@YourBotName`, or directly using the command. For example:    
<img width="400" alt="image" src="https://github.com/user-attachments/assets/c93196d9-8506-474b-8b09-1930b8bb42f1" />


All the above commands are also available in group chats without needing to switch to private chat.

