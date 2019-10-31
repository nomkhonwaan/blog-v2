import { OnInit, Component, Input } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-author',
  template: `
    <img src="assets/images/303589.webp" class="avatar" />
    <span *ngIf="displayName" class="display-name">Natcha Luangaroonchai</span>
  `,
  styles: [
    `
      :host {
        display: flex;
        justify-content: flex-start;
      }
    `,
    `
      .avatar {
        border-radius: 50%;
        height: 6.4rem;
        width: 6.4rem;
      }
    `,
    `
      .display-name {
        color: #333;
        display: inline-block;
        font: normal 400 1.6rem Lato, sans-serif;
        margin: .8rem 0 0 1.6rem;
      }
    `,
  ],
})
export class PostAuthorComponent extends PostComponent {

  @Input()
  displayName: boolean = true;

}
