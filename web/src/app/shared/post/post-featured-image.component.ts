import { Component, OnInit, Input, HostBinding, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-featured-image',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div *ngIf="coverMode; else nonCoverMode">
      <ng-content></ng-content>
    </div>

    <ng-template #nonCoverMode>
      <img *ngIf="src" [src]="src" class="lazyload">
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
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}?width=${this.innerWidth}&height=${this.innerHeight}`;
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
