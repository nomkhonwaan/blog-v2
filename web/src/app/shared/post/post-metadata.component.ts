import { OnInit, Component } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-metadata',
  template: `
    <ng-content></ng-content>
  `,
  styles: [],
})
export class PostMetadataComponent extends PostComponent implements OnInit {

  ngOnInit(): void {

  }

}
