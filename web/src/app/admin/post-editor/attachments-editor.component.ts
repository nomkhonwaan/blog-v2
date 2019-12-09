import { Component, Output, EventEmitter } from '@angular/core';
import { faImage, faSpinnerThird, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';

import { EditorComponent } from './editor.component';
import { forkJoin, Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'app-attachments-editor',
  template: `
    <div class="_section-title">
      {{'Attachments' | uppercase}}
    </div>

    <ul class="attachments _list-unstyled">
      <li class="attachment-item --new" [ngClass]="{'--disabled': isUploading}" (click)="isUploading ? false : file.click()">
        <fa-icon [icon]="isUploading ? faSpinnerThird : faImage" [spin]="isUploading"></fa-icon>

        <input #file type="file" name="files" multiple="multiple" (change)="onChange($event.target.files)" />
      </li>

      <li *ngFor="let attachment of post.attachments;let i = index" class="attachment-item" (click)="onClickAttachment(attachment)">
        <img src="/api/v2.1/storage/{{attachment.slug}}?width=64" alt="{{attachment.fileName}}" />
      </li>
    </ul>
    `,
  styles: [
    `
      .attachments {
        display: flex;
        flex-flow: row wrap;
        justify-content: space-between;
        margin: 0 3.2rem !important;
      }
    `,
    `
      .attachment-item {
        align-items: center;
        background: #eee;
        cursor: pointer;
        display: flex;
        justify-content: center;
        height: 4.2rem;
        margin-bottom: 0.8rem;
        width: 4.2rem;
      }
    `,
    `
      .attachment-item.--new {
        color: #bdbdbd;
        font-size: 2.2rem;
      }
    `,
    `
      .attachment-item.--disabled {
        cursor: default;
      }
    `,
    `
      .attachment-item > img {
        max-height: 100%;
        max-width: 100%;
      }
    `,
    `
      input[type="file"] {
        display: none;
      }
    `,
  ],
})
export class AttachmentsEditorComponent extends EditorComponent {

  faImage: IconDefinition = faImage;
  faSpinnerThird: IconDefinition = faSpinnerThird;
  faTimes: IconDefinition = faTimes;

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
