import { EditorComponent } from './editor.component';
import { Component } from '@angular/core';

@Component({
  selector: 'app-markdown-editor',
  template: `
    <textarea appAutoResize placeholder="Markdown" (change)="onChange()" [(ngModel)]="post.markdown"></textarea>
  `,
  styles: [
    `
      textarea {
        border: none;
        background: #fff;
        font: normal 300 1.6rem Pridi, sans-serif;
        min-height: calc(100% - 8.9rem - 3.2rem);
        /*                      title    margin-bottom */
        margin-top: 3.2rem;
        outline: none;
        overflow-y: hidden;
        padding: 3.2rem;
        width: 100%;
      }
    `,
    `
      textarea::placeholder {
        font: italic 400 1.6rem Lato, sans-serif;
      }
    `,
  ],
})
export class MarkdownEditorComponent extends EditorComponent {

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
