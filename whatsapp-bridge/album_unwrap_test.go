package main

import (
	"testing"
	"time"

	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestUnwrapAssociatedChildImage(t *testing.T) {
	inner := &waProto.Message{ImageMessage: &waProto.ImageMessage{
		URL: proto.String("https://mmg.whatsapp.net/x"), MediaKey: []byte{1, 2, 3},
	}}
	wrapped := &waProto.Message{AssociatedChildMessage: &waProto.FutureProofMessage{Message: inner}}

	if mt, _, _, _, _, _, _ := extractMediaInfo(wrapped, time.Now(), "ID"); mt != "" {
		t.Fatalf("precondition: wrapped message should hide the image, got %q", mt)
	}
	mt, _, url, key, _, _, _ := extractMediaInfo(unwrapExtraWrappers(wrapped), time.Now(), "ID")
	if mt != "image" || url == "" || len(key) == 0 {
		t.Fatalf("album child image not extracted: type=%q url=%q key=%v", mt, url, key)
	}
}

func TestUnwrapLeavesPlainMessage(t *testing.T) {
	m := &waProto.Message{Conversation: proto.String("hi")}
	if got := unwrapExtraWrappers(m); got != m {
		t.Fatal("plain message must be returned unchanged")
	}
}
