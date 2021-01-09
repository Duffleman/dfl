package pipeline

import (
	"dfl/lib/rpc"
)

var nsfwPage = rpc.MakeTemplate([]string{
	"./resources/short/nsfw.html",
	"./resources/short/layouts/root.html",
})

// HandleNSFWPrimer will show a NSFW primer screen if a resource has NSFW content
func HandleNSFWPrimer(p *Pipeline) (bool, error) {
	if p.context.forceDownload || p.context.primed {
		return true, nil
	}

	var anyNSFW bool

	for _, i := range p.rwqs {
		rwq := i

		if rwq.r.NSFW {
			anyNSFW = true
		}
	}

	if !anyNSFW {
		return true, nil
	}

	return false, nsfwPage.Execute(p.w, nil)
}
