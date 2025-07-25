package llm

import (
	"context"
	"errors"
	"fmt"
	"testing"
	
	openrouter "github.com/revrost/go-openrouter"
	"github.com/stretchr/testify/assert"
	"github.com/yincongcyincong/telegram-deepseek-bot/conf"
	"github.com/yincongcyincong/telegram-deepseek-bot/db"
	"github.com/yincongcyincong/telegram-deepseek-bot/param"
)

func TestOpenRouterSend(t *testing.T) {
	messageChan := make(chan *param.MsgInfo)
	
	go func() {
		for m := range messageChan {
			fmt.Println(m)
		}
	}()
	
	*conf.BaseConfInfo.Type = param.OpenRouter
	
	callLLM := NewLLM(WithChatId(1), WithMsgId(2), WithUserId("6"),
		WithMessageChan(messageChan), WithContent("hi"))
	callLLM.LLMClient.GetModel(callLLM)
	callLLM.LLMClient.GetMessages("6", "hi")
	err := callLLM.LLMClient.Send(context.Background(), callLLM)
	assert.Equal(t, nil, err)
	
}

func TestAIRouterReq_GetModel_Default(t *testing.T) {
	mock := &AIRouterReq{}
	l := &LLM{UserId: "non-exist"}
	
	mock.GetModel(l)
	assert.Equal(t, param.DeepseekDeepseekR1_0528Free, l.Model)
}

func TestAIRouterReq_GetModel_UserMode(t *testing.T) {
	
	mock := &AIRouterReq{}
	l := &LLM{UserId: "1"}
	mock.GetModel(l)
	assert.Equal(t, param.DeepseekDeepseekR1_0528Free, l.Model)
}

func TestAIRouterReq_GetMessages(t *testing.T) {
	userId := "user123"
	db.InsertMsgRecord(userId, &db.AQ{
		Question: "What is AI?",
		Answer:   "Artificial Intelligence",
		Content:  `[{"role":"tool","content":{"text":"Tool result"}}]`,
	}, true)
	
	r := &AIRouterReq{}
	r.GetMessages(userId, "Tell me more")
	assert.True(t, len(r.OpenRouterMsgs) >= 3)
	assert.Equal(t, "user", r.OpenRouterMsgs[0].Role)
	assert.Equal(t, "assistant", r.OpenRouterMsgs[2].Role)
	assert.Equal(t, "Tell me more", r.OpenRouterMsgs[len(r.OpenRouterMsgs)-1].Content.Multi[0].Text)
}

func TestAIRouterReq_GetUserMessage(t *testing.T) {
	mock := &AIRouterReq{}
	mock.GetUserMessage("Hello")
	assert.Equal(t, "user", mock.OpenRouterMsgs[0].Role)
	assert.Equal(t, "Hello", mock.OpenRouterMsgs[0].Content.Multi[0].Text)
}

func TestAIRouterReq_GetAssistantMessage(t *testing.T) {
	mock := &AIRouterReq{}
	mock.GetAssistantMessage("Hi!")
	assert.Equal(t, "assistant", mock.OpenRouterMsgs[0].Role)
	assert.Equal(t, "Hi!", mock.OpenRouterMsgs[0].Content.Multi[0].Text)
}

func TestAIRouterReq_AppendMessages(t *testing.T) {
	main := &AIRouterReq{
		OpenRouterMsgs: []openrouter.ChatCompletionMessage{{Role: "user", Content: openrouter.Content{Text: "Main"}}},
	}
	child := &AIRouterReq{
		OpenRouterMsgs: []openrouter.ChatCompletionMessage{{Role: "assistant", Content: openrouter.Content{Text: "Child"}}},
	}
	
	main.AppendMessages(child)
	assert.Len(t, main.OpenRouterMsgs, 2)
	assert.Equal(t, "Child", main.OpenRouterMsgs[1].Content.Text)
}

func TestAIRouterReq_GetMessage(t *testing.T) {
	mock := &AIRouterReq{}
	mock.GetMessage("user", "test")
	assert.Equal(t, "test", mock.OpenRouterMsgs[0].Content.Multi[0].Text)
	
	mock.GetMessage("assistant", "answer")
	assert.Equal(t, "answer", mock.OpenRouterMsgs[1].Content.Multi[0].Text)
}

func TestAIRouterReq_requestOneToolsCall_JSONError(t *testing.T) {
	r := &AIRouterReq{}
	calls := []openrouter.ToolCall{{
		ID: "call1",
		Function: openrouter.FunctionCall{
			Name:      "fakeTool",
			Arguments: "{invalid json}",
		},
	}}
	// should not panic
	r.requestOneToolsCall(context.Background(), calls)
	assert.Len(t, r.OpenRouterMsgs, 0)
}

func TestAIRouterReq_requestToolsCall_JSONError(t *testing.T) {
	call := openrouter.ChatCompletionStreamChoice{
		Delta: openrouter.ChatCompletionStreamChoiceDelta{
			ToolCalls: []openrouter.ToolCall{{
				ID:   "1",
				Type: "function",
				Function: openrouter.FunctionCall{
					Name:      "invalid",
					Arguments: "{"}}, // malformed JSON
			},
		},
	}
	r := &AIRouterReq{}
	err := r.requestToolsCall(context.Background(), call)
	assert.True(t, errors.Is(err, ToolsJsonErr))
}
