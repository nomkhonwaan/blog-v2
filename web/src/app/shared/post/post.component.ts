import { Input } from '@angular/core';

export abstract class PostComponent {

  @Input()
  post: Post;

}
