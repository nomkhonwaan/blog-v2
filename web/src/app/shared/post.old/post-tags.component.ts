import { Component } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-tags',
  template: `
    <ul class="tags _list-unstyled">
      <li *ngFor="let tag of post.tags">
        <a [routerLink]="['tag', tag.slug]">
          {{tag.name}}
        </a>
      </li>
    </ul>
  `,
  styles: [],
})
export class PostTagsComponent extends PostComponent { }
