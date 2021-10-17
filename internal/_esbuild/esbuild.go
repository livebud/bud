package esbuild

import (
	"fmt"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

type BuildOptions = esbuild.BuildOptions

type Plugin = esbuild.Plugin

var FormatESModule = esbuild.FormatESModule
var PlatformBrowser = esbuild.PlatformBrowser

func Build(options BuildOptions) (*esbuild.BuildResult, error) {
	result := esbuild.Build(options)
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Color:         true,
			Kind:          esbuild.ErrorMessage,
			TerminalWidth: 80,
		})
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}
	return &result, nil
}

// Plugin{
// 	Build(esbuild.BuildOptions{
// 	github.io/api/#how-conditions-work
// 	FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
// 	ErrorMessage,
// 	Plugin{
// 	Build(esbuild.BuildOptions{
// 	FormatESModule,
// 	PlatformBrowser,
// 	github.io/api/#how-conditions-work
// 	FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
// 	ErrorMessage,
// 	Plugin{
// 	EntryPoint, len(views))
// 	EntryPoint{
// 	Build(esbuild.BuildOptions{
// 	FormatESModule,
// 	PlatformBrowser,
// 	github.io/api/#how-conditions-work
// 	FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
// 	ErrorMessage,
// 	Plugin {
// 	Plugin{
// 	PluginBuild) {
// 	OnResolveOptions{Filter: `^bud\/view\/(?:[A-Za-z\-0-9]+\/)*_[A-Za-z\-0-9]+\.(svelte|jsx)$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
// 	OnLoadOptions{Filter: `.*`, Namespace: "dom"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
// 	LoaderJS
// 	Plugin {
// 	Plugin{
// 	PluginBuild) {
// 	OnResolveOptions{Filter: ".*"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
// 	Plugin {
// 	Plugin{
// 	PluginBuild) {
// 	OnLoadOptions{Filter: `\.svelte$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
// 	LoaderJS
