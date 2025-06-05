package main

import (
	"log"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp2(u User) error {

	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	go func() {
		for {
			if err := h.newsletterClient.AddToNewsletter(u); err == nil {
				return
			}
		}
	}()

	go func() {
		for {
			if err := h.notificationsClient.SendNotification(u); err == nil {
				return
			}
		}
	}()

	return nil
}

func (h Handler) SignUp(u User) error {

	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}
	log.Println("We managed to create the account")

	h.retry(u,
		h.newsletterClient.AddToNewsletter,
		h.notificationsClient.SendNotification,
	)

	return nil
}

func (h *Handler) retry(u User, fs ...func(u User) error) {
	for i, f := range fs {
		go func(index int, f func(u User) error) {
			defer log.Println("Done: ", index, "/", len(fs))
			for {
				if err := f(u); err == nil {
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
		}(i, f)
	}
}
