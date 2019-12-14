import { Component } from '@angular/core';

import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-title-editor',
  templateUrl: './title.component.html',
  styleUrls: ['./title.component.scss'],
})
export class PostTitleEditorComponent extends AbstractPostEditorComponent {

  onChange(): void {
    this.mutate(
      `
        mutation {
          updatePostTitle(slug: $slug, title: $title) {
            ...EditablePost
          }
        }
      `,
      {
        slug: this.post.slug,
        title: this.post.title,
      },
    );
  }

}
