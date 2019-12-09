import { Component } from '@angular/core';

import { EditorComponent } from './editor.component';

@Component({
  selector: 'app-title-editor',
  template: `
    <input placeholder="Title" (change)="onChange()" [(ngModel)]="post.title" />
  `,
  styles: [
    `
      input {
        border: none;
        background: #fff;
        font: normal 300 1.6rem Pridi, sans-serif;
        outline: none;
        padding: 3.2rem;
        width: 100%;
      }
    `,
    `
      input::placeholder {
        font: italic 400 1.6rem Lato, sans-serif;
      }
    `,
  ],
})
export class TitleEditorComponent extends EditorComponent {

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
