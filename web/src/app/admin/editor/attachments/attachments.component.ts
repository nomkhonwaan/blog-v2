import { Component, EventEmitter, Input, Output } from '@angular/core';
import { faImage, faSpinnerThird, IconDefinition } from '@nomkhonwaan/pro-light-svg-icons';
import { forkJoin, Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-attachments-editor',
  templateUrl: './attachments.component.html',
  styleUrls: ['./attachments.component.scss'],
})
export class PostAttachmentsEditorComponent extends AbstractPostEditorComponent {

  /**
   * A selected attachment
   */
  @Input()
  selectedAttachment: Attachment;

  /**
   * An attachment that has been selected from the list of attachments
   */
  @Output()
  selectAttachment: EventEmitter<Attachment> = new EventEmitter(null);

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faImage,
    faSpinnerThird,
  };

  /**
   * Use to display spinner while uploading attachments to the storage server
   */
  isUploading = false;

  onChange(files: FileList): void {
    this.isUploading = true;

    forkJoin(
      Array.
        from(files).
        map((file: File): Observable<Attachment> => this.api.uploadFile(file)),
    ).pipe(
      finalize((): void => { this.isUploading = false; }),
    ).subscribe((attachments: Attachment[]): void => {
      this.updatePostAttachments(this.post.attachments.concat(attachments));
    });
  }

  onClickAttachment(attachment: Attachment): void {
    this.selectAttachment.emit(attachment);
  }

  wasSelected(attachment: Attachment): boolean {
    return this.selectedAttachment && this.selectedAttachment.slug === attachment.slug;
  }

  private updatePostAttachments(attachments: Attachment[]): void {
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
        attachmentSlugs: attachments.map((attachment: Attachment): string => attachment.slug),
      }
    ).subscribe((post: Post): void => this.changeSuccess.emit(post));
  }

}
