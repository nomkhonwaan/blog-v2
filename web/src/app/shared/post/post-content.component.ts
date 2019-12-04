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
import { Store, select } from '@ngrx/store';
import 'fslightbox';

import { PostComponent } from './post.component';
import { isLightboxOpened, isLightboxClosed } from './post.actions';

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

  isLightboxClosed: boolean;

  @ViewChild('content', { static: false })
  private elementRef: ElementRef;

  constructor(
    @Inject(DOCUMENT) private document: Document,
    private http: HttpClient,
    private renderer: Renderer2,
    private store: Store<PostState>,
  ) {
    super();

    store
      .pipe(select('post', 'content', 'fslightbox', 'closed'))
      .subscribe((closed: boolean): void => {
        this.isLightboxClosed = closed;
      })
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
    this.addExtraQueryToImageSrc(imgs);
    this.renderGitHubGist(scripts);

    window.refreshFsLightbox();

    window.fsLightboxInstances[''].props.onOpen = (): void => {
      this.store.dispatch(isLightboxOpened());
    };

    window.fsLightboxInstances[''].props.onClose = (): void => {
      this.store.dispatch(isLightboxClosed());
    };
  }

  private addExtraClassNamesToAllImages(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      node.setAttribute('class', `${node.getAttribute('class')} lazyload`);
    });
  }

  private addExtraQueryToImageSrc(imgs: NodeList): void {
    imgs.forEach((node: Element): void => {
      const src: string = node.getAttribute('src')

      node.setAttribute('src', `${src}?width=${this.innerWidth}`);

      const anchor: Element = this.renderer.createElement('a');
      const img: Element = <Element>node.cloneNode();

      this.renderer.setAttribute(anchor, 'data-fslightbox', '');
      this.renderer.setAttribute(anchor, 'href', src);
      this.renderer.appendChild(anchor, img);
      this.renderImageCaption(img);

      node.replaceWith(anchor);
    });
  }

  private renderImageCaption(img: Element): void {
    const alt: string = img.getAttribute('alt');
    const caption: Element = this.renderer.createElement('div');

    this.renderer.addClass(caption, 'caption');
    this.renderer.appendChild(caption, this.renderer.createText(alt));

    img.insertAdjacentElement('afterend', caption);
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
