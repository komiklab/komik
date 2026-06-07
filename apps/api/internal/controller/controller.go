package controller

import (
	"context"

	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/component"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	components []component.Component
	cfg        *internal.Config
}

func NewController(cfg *internal.Config) *Controller {
	return &Controller{
		components: make([]component.Component, 0),
		cfg:        cfg,
	}
}

func (c *Controller) Init() {
	log.Debug().Msg("Init controller")
	for _, comp := range c.components {
		log.Debug().Msg("Init component: " + comp.GetName())
		comp.Init()
	}
}

func (c *Controller) Start() {
	log.Debug().Msg("Start controller")
	for _, comp := range c.components {
		log.Debug().Msg("Start component: " + comp.GetName())
		go comp.Start()
	}
}

func (c *Controller) Stop() {
	log.Debug().Msg("Stop controller")
	ctx := context.Background()
	for _, comp := range c.components {
		log.Debug().Msg("Stop component: " + comp.GetName())
		comp.Stop(ctx)
	}
}

func (c *Controller) AddComponent(comp component.Component) {
	c.components = append(c.components, comp)
}
