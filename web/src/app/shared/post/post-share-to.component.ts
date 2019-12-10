import { Component, Input, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { faFacebookF, faTwitter, IconDefinition } from '@fortawesome/free-brands-svg-icons';

import { PostComponent } from './post.component';

import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-post-share-to',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
      <span *ngIf="shareCount" class="share-count">
        {{shareCount}}
      </span>
      <a
        href="https://www.facebook.com/sharer/sharer.php?u={{getEncodedURL()}}&amp;src=sdkpreparse"
        class="share-to"
        target="_blank">
        <fa-icon class="icon-facebook" [icon]="faFacebookF"></fa-icon>
      </a>

      <a
        class="share-to"
        href="https://twitter.com/intent/tweet?text={{getURL()}}"
        target="_blank">
        <fa-icon class="icon-twitter" [icon]="faTwitter"></fa-icon>
      </a>
  `,
  styles: [
    `
      :host {
        color: #808080;
      }
    `,
    `
      .share-count {
        display: inline-block;
        font-size: 3.2rem;
        text-transform: uppercase;
      }
    `,
    `
      .share-to {
        align-items: center;
        border: .1rem solid #b3b3b3;
        border-radius: 50%;
        display: inline-flex;
        font-size: 2.2rem;
        height: 4.8rem;
        justify-content: center;
        margin: 0 1.2rem;
        width: 4.8rem;
      }
    `,
    `
      .share-to:first-child {
        margin-left: 0;
      }
    `,
    `
      .share-to:last-child {
        margin-right: 0;
      }
    `,
    `
      .share-to > .icon-facebook,
      .share-to > .icon-twitter {
        display: block;
      }
    `,
    `
      .share-to > .icon-facebook {
        margin-top: .3rem;
      }
    `,
    `
      .share-to > .icon-twitter {
        margin-top: .3rem;
        margin-left: .2rem;
      }
    `,
  ],
})
export class PostShareToComponent extends PostComponent implements OnInit {

  @Input()
  flow = 'row';

  /**
   * Use to sharing to the social network
   */
  url: string = environment.url;

  /**
   * Use to display number of social network engagement
   */
  shareCount: string;

  faFacebookF: IconDefinition = faFacebookF;
  faTwitter: IconDefinition = faTwitter;

  ngOnInit(): void {
    if (this.flow === 'column') {
      if (this.post.engagement.shareCount > 0) {
        this.shareCount = this.post.engagement.shareCount.toString();
      }
    }
  }

  getURL(): string {
    const publishedAt: Date = new Date(
      new Date(this.post.publishedAt)
        .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
    );

    return this.url + '/' + [
      publishedAt.getFullYear().toString(),
      (publishedAt.getMonth() + 1).toString(), // a month number start at 0 not 1
      publishedAt.getDate().toString(),
      this.post.slug,
    ].join('/');
  }

  getEncodedURL(): string {
    return encodeURIComponent(this.getURL());
  }
}
