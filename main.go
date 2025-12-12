package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/agentio/sidecar"
	"github.com/agentio/translate-sidecar/genproto/translatepb"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	var source string
	var target string
	var parent string
	var credentials string
	var address string
	var token string
	cmd := &cobra.Command{
		Use:   "translate TEXT",
		Short: "Translate with the Cloud Translation API",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := sidecar.NewClient(sidecar.ClientOptions{Address:address})
			if token != "" {
				client.Header.Set("authorization", "Bearer "+token)
			} else if credentials != "" {
				token, err := accessToken(cmd.Context(), credentials)
				if err != nil {
					return err
				}
				client.Header.Set("authorization", "Bearer "+token)
			}
			response, err := sidecar.CallUnary[translatepb.TranslateTextRequest, translatepb.TranslateTextResponse](
				cmd.Context(),
				client,
				"/google.cloud.translation.v3.TranslationService/TranslateText",
				sidecar.NewRequest(
					&translatepb.TranslateTextRequest{
						SourceLanguageCode: source,
						TargetLanguageCode: target,
						Contents:           args,
						Parent:             parent,
					}))
			if err != nil {
				return err
			}
			b, err := protojson.MarshalOptions{Indent: "  "}.Marshal(response.Msg)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "en-us", "source language")
	cmd.Flags().StringVar(&target, "target", "es-mx", "target language")
	cmd.Flags().StringVarP(&parent, "parent", "p", "", "parent project (format: projects/PROJECTID)")
	cmd.Flags().StringVar(&credentials, "credentials", "", "service account credentials")
	cmd.Flags().StringVarP(&address, "address", "a", "translate.googleapis.com:443", "service address")
	cmd.Flags().StringVar(&token, "token", "", "auth token")
	return cmd
}

func accessToken(ctx context.Context, credentials string) (string, error) {
	serviceAccountJSON, err := os.ReadFile(credentials)
	if err != nil {
		return "", err
	}
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	creds, err := google.CredentialsFromJSON(ctx, serviceAccountJSON, scopes...)
	if err != nil {
		return "", err
	}
	if creds == nil {
		return "", errors.New("no credentials")
	}
	if creds.TokenSource == nil {
		return "", errors.New("no token source")
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}
	if token == nil {
		return "", errors.New("no token")
	}
	return token.AccessToken, nil
}
