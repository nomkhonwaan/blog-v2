import { Component, OnInit } from '@angular/core';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';
import { BehaviorSubject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: 'app-post-markdown-editor',
  templateUrl: './markdown.component.html',
  styleUrls: ['./markdown.component.scss'],
})
export class PostMarkdownEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * Use for debouncing keypress event on markdown value
   */
  markdown$: BehaviorSubject<string> = new BehaviorSubject('');

  ngOnInit(): void {
    this.markdown$.next(this.post.markdown);

    this.markdown$.pipe(debounceTime(1600)).subscribe((markdown: string): void => {
      this.updatePostContent(this.post.slug, markdown);
    });
  }

  onChange(): void {
    this.updatePostContent(this.post.slug, this.post.markdown);
  }

  onKeyPress(): void {
    this.markdown$.next(this.post.markdown);
  }

  private updatePostContent(slug: string, markdown: string): void {
    this.mutate(
      `
        mutation {
          updatePostContent(slug: $slug, markdown: $markdown) {
            ...EditablePost
          }
        }
      `,
      {
        slug,
        markdown,
      }
    );
  }
}
