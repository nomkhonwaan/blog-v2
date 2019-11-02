import { Component, Input } from '@angular/core';
import { faFacebookF, faTwitter, IconDefinition } from '@fortawesome/free-brands-svg-icons';
import { faCopy, IconDefinition as SolidIconDefinition } from '@fortawesome/fontawesome-free-solid';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-share-to',
  template: `
      <a class="icon" href="" target="_blank">
        <fa-icon [icon]="faFacebookF"></fa-icon>
      </a>

      <a class="icon" href="" target="_blank">
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
        width: 4.8rem;
      }
    `,
  ],
})
export class PostShareToComponent extends PostComponent {

  @Input()
  flow: string = 'row';

  faFacebookF: IconDefinition = faFacebookF;
  faTwitter: IconDefinition = faTwitter;
  faCopy: SolidIconDefinition = faCopy;

}
