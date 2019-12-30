import { Component, EventEmitter, Input, Output } from '@angular/core';
import { faSpinnerThird, faTrash, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import update from 'immutability-helper';
import { finalize } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-attachment-viewer',
  templateUrl: './attachment-viewer.component.html',
  styleUrls: ['./attachment-viewer.component.scss'],
})
export class PostAttachmentViewerComponent extends AbstractPostEditorComponent {

  /**
   * An attachment item
   */
  @Input()
  attachment: Attachment;

  /**
   * For emitting on view closed event
   */
  @Output()
  closed: EventEmitter<null> = new EventEmitter(null);

  /**
   * Use to display spinner while deleting attachment from the storage server
   */
  isDeleting = false;

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faTrash,
    faSpinnerThird,
  };

  close(): void {
    if (!this.isDeleting) {
      this.closed.emit(null);
    }
  }

  deleteFile(): void {
    if (confirm('Are you sure? this cannot be undone')) {
      this.isDeleting = true;
      const i: number = this.post.attachments.findIndex(({ slug }): boolean => slug === this.attachment.slug);

      this.api.deleteFile(this.attachment.slug).subscribe((): void => {
        this.mutate(
          `
          mutation {
            updatePostAttachments(slug: $slug, attachmentSlugs: $attachmentSlugs) {
              ...EditablePost
            }
          }
        `,
          {
            slug: this.post.slug,
            attachmentSlugs: update(this.post, { attachments: { $splice: [[i, 1]] } }).attachments
              .map((attachment: Attachment) => attachment.slug),
          },
        ).pipe(
          finalize((): void => {
            this.isDeleting = false;
            this.close();
          }),
        ).subscribe((post: Post): void => this.changeSuccess.emit(post));
      });
    }
  }

}
