package main

import (
	"fmt"

	"github.com/keepcalmist/chat-service/internal/config"
	jobsrepo "github.com/keepcalmist/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	msgproducer "github.com/keepcalmist/chat-service/internal/services/msg-producer"
	"github.com/keepcalmist/chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-sent"
	sendclientmessagejob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/keepcalmist/chat-service/internal/store"
)

func initOutbox(
	cfg config.Services,
	database *store.Database,
	repoJobs *jobsrepo.Repo,
	repoMsg *messagesrepo.Repo,
	producer *msgproducer.Service,
	stream eventstream.EventStream,
) (*outbox.Service, error) {
	outboxService, err := outbox.New(outbox.NewOptions(
		cfg.Outbox.Workers,
		cfg.Outbox.IdleTime,
		cfg.Outbox.ReserveFor,
		repoJobs,
		database,
	))
	if err != nil {
		return nil, fmt.Errorf("init outbox service: %v", err)
	}

	sendClientMsgJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(repoMsg, producer, stream))
	if err != nil {
		return nil, fmt.Errorf("init send client message job: %v", err)
	}

	clientMessageSentJob, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(repoMsg, stream))
	if err != nil {
		return nil, fmt.Errorf("init client message sent job: %v", err)
	}

	clientMessageBlockedJob, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(repoMsg, stream))
	if err != nil {
		return nil, fmt.Errorf("init client message blocked job: %v", err)
	}

	err = outboxService.RegisterJob(sendClientMsgJob)
	if err != nil {
		return nil, fmt.Errorf("register send client message job: %v", err)
	}

	err = outboxService.RegisterJob(clientMessageSentJob)
	if err != nil {
		return nil, fmt.Errorf("register client send message job: %v", err)
	}

	err = outboxService.RegisterJob(clientMessageBlockedJob)
	if err != nil {
		return nil, fmt.Errorf("register client blocked message job: %v", err)
	}

	return outboxService, nil
}
