import { DOCUMENT } from '@angular/common';
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  OnInit,
  Renderer2,
  ViewChild,
  Inject,
} from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-content',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <article #content></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit, AfterViewInit {

  content: string;

  @ViewChild('content', { static: false })
  private elementRef: ElementRef;

  constructor(@Inject(DOCUMENT) private document: Document, private http: HttpClient, private renderer: Renderer2) {
    super();
  }

  ngOnInit(): void {
    this.content = this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments');
  }

  ngAfterViewInit(): void {
    this.elementRef.nativeElement.appendChild(
      this.document.createRange().createContextualFragment(this.content),
    );

    const imgs: NodeList = this.elementRef.nativeElement.querySelectorAll('img');
    const scripts: NodeList = this.elementRef.nativeElement.querySelectorAll('script');

    this.addExtraClassNamesToAllImages(imgs);
    this.renderAllImageCaptions(imgs);
    this.renderGitHubGist(scripts);
  }

  private addExtraClassNamesToAllImages(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      node.setAttribute('class', `${node.getAttribute('class')} lazyload`);
    });
  }

  private renderAllImageCaptions(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      const alt: string = node.getAttribute('alt');
      const caption: Element = this.renderer.createElement('div');

      this.renderer.addClass(caption, 'caption');
      this.renderer.appendChild(caption, this.renderer.createText(alt));

      node.insertAdjacentElement('afterend', caption);
    });
  }

  private renderGitHubGist(scripts: NodeList): void {
    scripts.forEach((node: Element): void => {
      const src: string = node.getAttribute('src');

      this.http.get(`/api/v2.1/github/gist?src=${encodeURIComponent(src)}`)
        .subscribe((res: Gist): void => {
          node.insertAdjacentHTML('afterend', res.div);

          const link: Element = this.renderer.createElement('link');
          link.setAttribute('rel', 'stylesheet');
          link.setAttribute('href', res.stylesheet);
          node.insertAdjacentElement('afterend', link);
        });
    });
  }

}
