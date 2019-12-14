import { Component } from '@angular/core';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-markdown-editor',
  templateUrl: './markdown.component.html',
  styleUrls: ['./markdown.component.scss'],
})
export class PostMarkdownEditorComponent extends AbstractPostEditorComponent {

  onChange(): void {
    this.mutate(
      `
        mutation {
          updatePostContent(slug: $slug, markdown: $markdown) {
            ...EditablePost
          }
        }
      `,
      {
        slug: this.post.slug,
        markdown: this.post.markdown,
      }
    );
  }

}
