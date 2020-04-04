import { DOCUMENT } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { AfterViewInit, ChangeDetectionStrategy, Component, ElementRef, Inject, OnInit, Renderer2, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { select, Store } from '@ngrx/store';
import 'fslightbox';
import { AbstractPostComponent } from '../abstract-post.component';
import { isLightboxClosed, isLightboxOpened } from '../post.actions';

@Component({
  selector: 'app-post-content',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './content.component.html',
  styleUrls: ['./content.component.scss'],
})
export class PostContentComponent extends AbstractPostComponent implements OnInit, AfterViewInit {

  /**
   * An HTML content of the post
   */
  html: string;

  /**
   * Use to indicate whether Lightbox has open or not
   */
  isLightboxClosed: boolean;

  /**
   * For escaping XSS protection while rendering HTML with Angular
   */
  @ViewChild('content')
  private content: ElementRef;

  constructor(
    @Inject(DOCUMENT) private document: Document,
    private http: HttpClient,
    private renderer: Renderer2,
    private router: Router,
    private store: Store<{ post: PostState }>,
  ) {
    super();
  }

  ngOnInit(): void {
    this.store
      .pipe(select('post', 'content', 'fslightbox', 'closed'))
      .subscribe((closed: boolean): void => {
        this.isLightboxClosed = closed;
      });
  }

  ngAfterViewInit(): void {
    this.content.nativeElement.appendChild(
      this.document.createRange().createContextualFragment(this.post.html),
    );

    const anchors: NodeList = this.content.nativeElement.querySelectorAll('a');
    const imgs: NodeList = this.content.nativeElement.querySelectorAll('img');
    const scripts: NodeList = this.content.nativeElement.querySelectorAll('script');

    this.fixHashLinks(anchors);
    this.addExtraClassNamesToAllImages(imgs);
    this.addExtraQueryToImageSrc(imgs);
    this.renderGitHubGists(scripts);

    if (imgs.length > 0) {
      window.refreshFsLightbox();

      window.fsLightboxInstances[''].props.onOpen = (): void => {
        this.store.dispatch(isLightboxOpened());
      };

      window.fsLightboxInstances[''].props.onClose = (): void => {
        this.store.dispatch(isLightboxClosed());
      };
    }
  }

  private fixHashLinks(anchors: NodeList): void {
    anchors.forEach((node: Element): void => {
      const href: string = node.getAttribute('href');

      if (/^#/.test(href)) {
        node.setAttribute('href', this.router.url + href);
      }
    });
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
      const img: Element = node.cloneNode() as Element;

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

  private renderGitHubGists(scripts: NodeList): void {
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
