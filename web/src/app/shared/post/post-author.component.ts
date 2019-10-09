import { OnInit, Component } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-author',
  template: `
    <ng-content></ng-content>
  `,
  styles: [],
})
export class PostAuthorComponent extends PostComponent implements OnInit {

  ngOnInit(): void {

  }

}
