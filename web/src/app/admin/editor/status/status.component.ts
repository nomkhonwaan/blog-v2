import { Component, OnInit } from '@angular/core';
import update from 'immutability-helper';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-status-editor',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.scss'],
})
export class PostStatusEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * Available actions
   */
  availableStatuses: Array<DropdownItem> = [{ label: 'Draft', value: 'DRAFT' }, { label: 'Publish', value: 'PUBLISHED' }];

  /**
   * Use to display a current post status
   */
  currentStatus: DropdownItem;

  ngOnInit(): void {
    this.currentStatus = this.availableStatuses.
      map((item: DropdownItem): DropdownItem | null =>
        (item.value || item.label).toString() === this.post.status
          ? update(item, { label: { $set: item.label.toUpperCase() } })
          : null,
      ).
      filter((item: DropdownItem): DropdownItem => item)[0];

      console.log(this.currentStatus, this.post);
  }

  onChange(selectedItem: { label: string, value?: any }): void {
    this.mutate(
      `
          mutation {
            updatePostStatus(slug: $slug, status: $status) {
              ...EditablePost
            }
          }
        `,
      {
        slug: this.post.slug,
        status: selectedItem.value.toString(),
      },
    );

    this.currentStatus = update(selectedItem, { label: { $set: selectedItem.label.toUpperCase() } });
  }

}
