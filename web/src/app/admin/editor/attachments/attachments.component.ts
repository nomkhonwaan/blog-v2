import { Component, EventEmitter, Output } from '@angular/core';
import { faImage, faSpinnerThird, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
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
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faImage,
    faSpinnerThird,
  };

  /**
   * An attachment that has been selected from the list of attachments
   */
  @Output()
  selectAttachment: EventEmitter<Attachment> = new EventEmitter(null);

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
    );
  }

}
