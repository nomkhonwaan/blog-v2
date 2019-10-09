import { OnInit, Component } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-author',
  template: `
    <img src="assets/images/303589.png" class="avatar" />
    <div class="display-name">Natcha Luangaroonchai</div>
  `,
  styles: [
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
          font: normal 400 1.6rem Lato, sans-serif;
      }
    `,
  ],
})
export class PostAuthorComponent extends PostComponent implements OnInit {

  ngOnInit(): void {

  }

}
