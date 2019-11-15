import { PostComponent } from './post.component';
import { Component, OnInit, Input, HostBinding } from '@angular/core';

@Component({
  selector: 'app-post-featured-image',
  template: `
    <div *ngIf="coverMode; else nonCoverMode">
      <ng-content></ng-content>
    </div>

    <ng-template #nonCoverMode>
      <img *ngIf="src" [src]="src">
    </ng-template>
  `,
  styleUrls: ['./post-featured-image.component.scss'],
})
export class PostFeaturedImageComponent extends PostComponent implements OnInit {

  /**
   * Used to indicate whether featured image should display as cover or not
   */
  @Input()
  @HostBinding('class.-cover-mode')
  coverMode = false;

  @HostBinding('class.-with-featured-image')
  withFeaturedImage = false;

  @HostBinding('style.background-image')
  src: string;

  ngOnInit(): void {
    if (this.hasFeaturedImage()) {
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}`;
      this.withFeaturedImage = true;

      if (this.coverMode) {
        this.src = `url(${this.src})`;
      }
    }
  }

  hasFeaturedImage(): boolean {
    return !!this.post && this.post.featuredImage.slug.length > 0;
  }

}
