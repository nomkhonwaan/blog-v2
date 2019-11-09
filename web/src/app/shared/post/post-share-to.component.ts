import { Component, Input } from '@angular/core';
import { faFacebookF, faTwitter, IconDefinition } from '@fortawesome/free-brands-svg-icons';
import { faCopy, IconDefinition as SolidIconDefinition } from '@fortawesome/fontawesome-free-solid';

import { PostComponent } from './post.component';

import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-post-share-to',
  template: `
      <a
        href="https://www.facebook.com/sharer/sharer.php?u={{getEncodedURL()}}&amp;src=sdkpreparse"
        class="icon"
        target="_blank">
        <fa-icon [icon]="faFacebookF"></fa-icon>
      </a>

      <a
        class="icon"
        href="https://twitter.com/intent/tweet?text={{getURL()}}"
        target="_blank">
        <fa-icon [icon]="faTwitter"></fa-icon>
      </a>
  `,
  styles: [
    `
      .icon {
        align-items: center;
        border: .1rem solid #b3b3b3;
        border-radius: 50%;
        color: #808080;
        display: inline-flex;
        font-size: 2.2rem;
        height: 4.8rem;
        justify-content: center;
        margin: 0 1.2rem;
        width: 4.8rem;
      }
    `,
    `
      .icon:first-child {
        margin-left: 0;
      }
    `,
    `
      .icon:last-child {
        margin-right: 0;
      }
    `,
  ],
})
export class PostShareToComponent extends PostComponent {

  @Input()
  flow = 'row';

  /**
   * Used to sharing to the social network
   */
  url: string = environment.url;

  faFacebookF: IconDefinition = faFacebookF;
  faTwitter: IconDefinition = faTwitter;
  faCopy: SolidIconDefinition = faCopy;

  getURL(): string {
    const publishedAt: Date = new Date(
      new Date('2016-04-01T17:00:00Z')
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
