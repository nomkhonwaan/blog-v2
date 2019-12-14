import { DOCUMENT } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { AfterViewInit, ChangeDetectionStrategy, Component, ElementRef, Inject, OnInit, Renderer2, ViewChild } from '@angular/core';
import { select, Store } from '@ngrx/store';
import 'fslightbox';
import { isLightboxClosed, isLightboxOpened } from '../post.actions';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-content',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './content.component.html',
  styleUrls: ['./content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit, AfterViewInit {

  /**
   * An HTML content of the post
   */
  content: string;

  /**
   * Use to indicate whether Lightbox has open or not
   */
  isLightboxClosed: boolean;

  /**
   * For escaping XSS protection while rendering HTML with Angular
   */
  @ViewChild('content', { static: false })
  private elementRef: ElementRef;

  constructor(
    @Inject(DOCUMENT) private document: Document,
    private http: HttpClient,
    private renderer: Renderer2,
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
    console.log(this.content);
    this.elementRef.nativeElement.appendChild(
      this.document.createRange().createContextualFragment(this.post.html),
    );

    const imgs: NodeList = this.elementRef.nativeElement.querySelectorAll('img');
    const scripts: NodeList = this.elementRef.nativeElement.querySelectorAll('script');

    this.addExtraClassNamesToAllImages(imgs);
    this.addExtraQueryToImageSrc(imgs);
    this.renderGitHubGist(scripts);

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
