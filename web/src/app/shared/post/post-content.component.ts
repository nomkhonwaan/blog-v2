import { Component, OnInit, Directive, ElementRef, AfterViewInit, Renderer2 } from '@angular/core';

import { PostComponent } from './post.component';

@Directive({
  selector: '[appPostContent]',
})
export class PostContentDirective implements AfterViewInit {

  constructor(private el: ElementRef, private renderer: Renderer2) { }

  ngAfterViewInit(): void {
    this.renderImageCaptionFromItsAltAttribute();
  }

  renderImageCaptionFromItsAltAttribute(): void {
    const imgs: NodeList = this.el.nativeElement.querySelectorAll('img[alt]');

    imgs.forEach((node: Element): void => {
      const alt: string = node.getAttribute('alt');
      const caption: Element = this.renderer.createElement('div');

      this.renderer.addClass(caption, 'caption');
      this.renderer.appendChild(caption, this.renderer.createText(alt));

      node.insertAdjacentElement('afterend', caption);
    });
  }
}

@Component({
  selector: 'app-post-content',
  template: `
    <article appPostContent [innerHTML]="content"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: string;

  ngOnInit(): void {
    this.content = this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments');
  }
}
