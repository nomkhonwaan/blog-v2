import { Input } from '@angular/core';

import { Type } from './type.enum';

export abstract class PostComponent {

  @Input()
  post: Post;

  @Input()
  type: Type;

  single: Type = Type.Single;
  medium: Type = Type.Medium;
  thumbnail: Type = Type.Thumbnail;

}
