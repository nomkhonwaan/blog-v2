import { Component, OnInit, Directive, ElementRef, AfterViewInit, Renderer2, ChangeDetectionStrategy, Input } from '@angular/core';

import { PostComponent } from './post.component';

@Directive({
  selector: '[appPostContent]',
})
export class PostContentDirective implements AfterViewInit {

  @Input()
  innerWidth: number;

  @Input()
  innerHeight: number;

  constructor(private el: ElementRef, private renderer: Renderer2) { }

  ngAfterViewInit(): void {
    const imgs: NodeList = this.el.nativeElement.querySelectorAll('img');

    this.addExtraClassToImageClass(imgs);
    this.addExtraQueryToImageSrc(imgs);
    this.renderImageCaption(imgs);
  }

  addExtraClassToImageClass(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      const classes: string = node.getAttribute('class');

      node.setAttribute('class', `${classes} lazyload`);
    })
  }

  addExtraQueryToImageSrc(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      const src: string = node.getAttribute('src')

      node.setAttribute('src', `${src}?width=${this.innerWidth}&height=${this.innerHeight}`);
    });
  }

  renderImageCaption(imgs: NodeList): void {
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
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <article appPostContent [innerHTML]="content" [innerWidth]="innerWidth" [innerHeight]="innerHeight"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: string;

  ngOnInit(): void {
    this.content = this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments');
  }
}
