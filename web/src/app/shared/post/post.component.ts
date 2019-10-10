import { Input, Component } from '@angular/core';

export abstract class PostComponent {

  @Input()
  post: Post;

}
