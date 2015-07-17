// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package steven

import (
	"fmt"
	"io"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
)

var progressBars []*progressBar

type progressBar struct {
	id   int
	bar  *ui.Image
	text *ui.Text
}

func newProgressBar() *progressBar {
	p := &progressBar{
		id: len(progressBars),
	}
	p.bar = ui.NewImage(render.GetTexture("solid"), 0, 21*float64(p.id), 854, 21, 0, 0, 1, 1, 0, 125, 0)
	ui.AddDrawable(p.bar.Attach(ui.Top, ui.Left))
	p.text = ui.NewText("", 1, 21*float64(p.id)+1, 255, 255, 255)
	ui.AddDrawable(p.text.Attach(ui.Top, ui.Left))
	progressBars = append(progressBars, p)
	p.bar.SetLayer(-241)
	p.text.SetLayer(-240)
	return p
}

func (p *progressBar) update(progress float64, text string) {
	p.text.Update(text)
	width, _ := window.GetFramebufferSize()
	sw := 854 / float64(width)
	if ui.DrawMode == ui.Unscaled {
		sw = ui.Scale
		p.bar.SetWidth((854 / sw) * progress)
	} else {
		p.bar.SetWidth(float64(width) * progress)
	}
}

func (p *progressBar) remove() {
	p.bar.Remove()
	p.text.Remove()
	for i, po := range progressBars {
		if po == p {
			progressBars = append(progressBars[:i], progressBars[i+1:]...)
			break
		}
	}
	for i, po := range progressBars {
		po.id = i
		po.bar.SetY(21 * float64(po.id))
		po.text.SetY(21*float64(po.id) + 1)
	}
}

func (p *progressBar) watchCopy(dst io.Writer, src io.Reader, size int64, text string) (n int64, err error) {
	buf := make([]byte, 4906)
	for {
		nn, err := src.Read(buf)
		if nn > 0 {
			_, err := dst.Write(buf[:nn])
			n += int64(nn)
			progress := float64(n) / float64(size)
			syncChan <- func() {
				if size == -1 {
					p.update(0.5, fmt.Sprintf(text, "??"))
				} else {
					p.update(progress, fmt.Sprintf(text, int(100*progress)))
				}
			}
			if err != nil {
				return n, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}
	}
	return
}
