/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Modifications made
 1. Pass around release.Release pointers rather than interfaces.
 2. Removed many handlers.
 3. Added the googlechat handler.
*/

package handlers

import (
	"github.com/larderdev/kubewise/handlers/googlechat"
	"github.com/larderdev/kubewise/handlers/slack"
	"github.com/larderdev/kubewise/kwrelease"
)

type Handler interface {
	Init() error
	HandleEvent(releaseEvent *kwrelease.Event)
}

var Map = map[string]interface{}{
	"default":    &Default{},
	"slack":      &slack.Slack{},
	"googlechat": &googlechat.GoogleChat{},
}

type Default struct {
}

func (d *Default) Init() error {
	return nil
}

func (d *Default) HandleEvent(releaseEvent *kwrelease.Event) {

}
