import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  Directive,
  ElementRef,
  Input,
  OnInit,
  Renderer2,
} from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';

import { PostComponent } from './post.component';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { environment } from 'src/environments/environment';

@Directive({
  selector: '[appPostContent]',
})
export class PostContentDirective implements AfterViewInit {

  @Input()
  innerWidth: number;

  @Input()
  innerHeight: number;

  constructor(private el: ElementRef, private http: HttpClient, private renderer: Renderer2) { }

  ngAfterViewInit(): void {
    const imgs: NodeList = this.el.nativeElement.querySelectorAll('img');
    const scripts: NodeList = this.el.nativeElement.querySelectorAll('script');

    this.addExtraClassToImageClass(imgs);
    this.addExtraQueryToImageSrc(imgs);
    this.renderImageCaption(imgs);
    this.renderGist(scripts);
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

  renderGist(scripts: NodeList): void {
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

@Component({
  selector: 'app-post-content',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <article appPostContent [innerHTML]="content" [innerWidth]="innerWidth" [innerHeight]="innerHeight"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: SafeHtml;

  constructor(private sanitizer: DomSanitizer) {
    super();
  }

  ngOnInit(): void {
    this.content = this.sanitizer.bypassSecurityTrustHtml(
      this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments'),
    );
  }

}
