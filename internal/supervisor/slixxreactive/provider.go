package slixxreactive

import (
	"context"
	"github.com/samsarahq/thunder/reactive"
)

type SlixxReactive struct {
	cause    string
	resource *reactive.Resource
}

var reactives = make([]*SlixxReactive, 0)

func InvalidateOn(ctx context.Context, cause string) {
	r := reactive.NewResource()
	reactive.AddDependency(ctx, r, nil)

	reactive := &SlixxReactive{
		cause:    cause,
		resource: r,
	}
	reactives = append(reactives, reactive)
}

func Event(cause string) {
	var reactivesToRemove = make([]*SlixxReactive, 0)
	for _, reactiveModel := range reactives {
		if reactiveModel.cause == cause {
			reactiveModel.resource.Invalidate()
			reactivesToRemove = append(reactivesToRemove, reactiveModel)
		}
	}
	reactives = remove(reactives, reactivesToRemove)
}

func remove(slixxReactives []*SlixxReactive, remove []*SlixxReactive) []*SlixxReactive {
	for _, reactiveToRemove := range remove {
		for i, reactive := range slixxReactives {
			if reactive == reactiveToRemove {
				slixxReactives = append(slixxReactives[:i], slixxReactives[i+1:]...)
				break
			}
		}
	}
	return slixxReactives
}
