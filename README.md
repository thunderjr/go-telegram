# Go Telegram

This repository features a library that offers a more intuitive API for implementing bot logic on Telegram, built on top of the `tgbotapi` library.
The library simplifies the process of handling keyboard callbacks, webapp data and form submissions by using the Finite State Machine pattern for form management.

### Globals and Data Providers

The library includes a `globals` configuration file that defines key settings and data providers for the bot's functionality. Currently, it supports Redis as the data provider through the `RedisProvider` option.

#### Key Global Variables:
- **`AppName`**: This variable should be set to the name of your application. It is required to initialize the repository.
- **`DataProvider`**: Specifies the data provider to use. The default and currently only supported provider is Redis (`RedisProvider`).

#### Functions:
- **`SetAppName(name string)`**: Sets the application name, which is used as a prefix in the Redis keys.
- **`SetDataProvider(provider dataProvider)`**: Allows switching the data provider, though Redis is the only implementation at this time.

#### Usage Example:

```go
package main

import (
	"context"
	"github.com/thunderjr/go-telegram/pkg/bot"
)

func main() {
	bot.SetAppName("my_telegram_bot")
	bot.SetDataProvider(bot.RedisProvider)

	repo := bot.NewRepository[MyEntity]()
	// Use the repository as needed
}
```

#### Repository Context Functions:
- **`WithEditableRepo`**: Adds an editable message repository to the context.
- **`WithReplyActionRepo`**: Adds a reply action repository to the context.
- **`WithFormAnswerRepo`**: Adds a form answer repository to the context.

## Examples

### Simple Message

```go
package main

import (
	"context"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

func main() {
	bot, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	msg := message.NewSimpleMessage(&message.Params{
		Content:   "Hello, this is a simple message!",
		Recipient: 42069,
		Bot:       bot,
	})

	if _, err := msg.Send(context.TODO(),
		// Parse mode
		message.WithMarkdownParseMode(),
		// Buttons (check handlers for each callbacks type)
		message.WithMessageButtons(message.KeyboardRow{{"Click Me", "simple"}}),
		// OR
		message.WithWebappButton("Google", "https://google.com"),
		// Reply options
		message.WithReplyToMessageID(69420),
		// OR
		message.WithForceReply(),
	); err != nil {
		panic(err)
	}
}
```

### Keyboard Callback

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

func main() {
	updateHandlers := []update.Handler{
		update.NewKeyboardCallbackUpdate("simple", func(u tgbotapi.Update) error {
			fmt.Println("Button clicked:", u.CallbackQuery.Message.Text)
			return nil
		}),
	}

	bot, err := bot.New(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		bot.WithUpdateHandlers(updateHandlers),
	)
	if err != nil {
		panic(err)
	}

	msg := message.NewSimpleMessage(&message.Params{
		Content:   "Hello, this is a simple message!",
		Recipient: 42069,
		Bot:       bot,
	})

	if _, err := msg.Send(context.TODO(),
		message.WithMessageButtons(message.KeyboardRow{{"Click Me", "simple"}}),
	); err != nil {
		panic(err)
	}

	chErr := make(chan error)
	go func() {
		for err := range chErr {
			log.Printf("error: %+v\n", err)
		}
	}()

	bot.Updates(context.TODO(), chErr)
}
```

### Sending and Editing a Messages

This example demonstrates sending an `EditableMessage` that can be updated later with a `CandidateMessage`.

```go
package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
	"github.com/thunderjr/go-telegram/pkg/bot"
)

func main() {
	bot, err := bot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	simpleMsg := message.NewSimpleMessage(&message.Params{
		Content:   "Hello, this is a simple message!",
		Recipient: 42069,
		Bot:       bot,
	})

	// Convert the simple message to editable
	editableMsg := message.ToEditable(simpleMsg)
	if _, err := editableMsg.Send(ctx); err != nil {
		panic(err)
	}

	candidateMsg := message.NewCandidateMessage(&message.Params{
		Content:   "This is a candidate message that can be used to edit the last editable message.",
		Recipient: 42069,
		Bot:       bot,
	})

	if _, err := candidateMsg.Send(ctx); err != nil {
		panic(err)
	}
}
```

### Handling a Reply Action

```go
package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

func main() {
	bot, err := message.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx = bot.WithReplyActionRepo(ctx, bot.NewRepository[message.ReplyAction]())

	msg := message.NewSimpleMessage(&message.Params{
		Content:   "Please reply to this message.",
		OnReply:   "handle_reply_action",
		Recipient: 42069,
		Bot:       bot,
	})

	if _, err := msg.Send(ctx); err != nil {
		panic(err)
	}

	// Function to handle the reply action
	handleReplyAction := func(ctx context.Context, msg *tgbotapi.Message) {
		log.Printf("Received reply: %s", msg.Text)
		// Perform additional logic based on the reply
	}

	bot.AddHandlers(
		update.NewReplyUpdate("handle_reply_action", handleReplyAction),
	)
}
```

### Form Submission

This example demonstrates how to handle form submissions using the Finite State Machine pattern.
You can also handle form fields with custom prompts and orders using the `PromptProvider` interface.
Implement the `FieldPrompts` method to provide custom prompts for each field.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

var _ update.PromptProvider = (*UserForm)(nil)

type UserForm struct {
	FirstName string `telegram_prompt:"What is your first name?" telegram_prompt_order:"1"`
	City      string `telegram_prompt:"Where do you live?"`
	Age       int
}

func (f *UserForm) FieldPrompts() ([]update.FormFieldPrompt, error) {
	return []update.FormFieldPrompt{
		{
			Name:   "Age",
			Prompt: "How old are you?",
			Order:  2,
		},
	}, nil
}

func main() {
	ctx := context.Background()
	ctx = bot.WithFormAnswerRepo(ctx, bot.NewRepository[update.FormAnswer]())

	bot, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	msg := message.NewSimpleMessage(&message.Params{
		Content:   "Please fill out the form.",
		Recipient: 42069,
		Bot:       bot,
	})

	if _, err := msg.Send(ctx,
		message.WithMessageButtons(message.KeyboardRow{{"Start Form", "start_form"}}), // Key should match the handler key to trigger the form
	); err != nil {
		panic(err)
	}

	onSubmit := func(ctx context.Context, data *UserForm) error {
		fmt.Println("Form submitted:", data)
		return nil
	}

	bot.AddHandlers(
		update.NewFormHandlers(ctx, &update.NewFormHandlerParams[UserForm]{
			Type:     update.HandlerTypeKeyboardCallback,
			Key:      "start_form", // Keyboard callback key
			Form:     &UserForm{},
			OnSubmit: onSubmit,
			Bot:      bot,
		})...,
	)

	chErr := make(chan error)
	go func() {
		for err := range chErr {
			log.Printf("error: %+v\n", err)
		}
	}()

	bot.Updates(ctx, chErr)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.
