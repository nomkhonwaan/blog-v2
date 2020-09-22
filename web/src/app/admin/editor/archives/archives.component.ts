import { Component, OnInit, Input } from '@angular/core';
import { faSearch, IconDefinition } from '@nomkhonwaan/pro-light-svg-icons';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

type Archive = Category | Tag;

@Component({
  selector: 'app-post-archives-editor',
  templateUrl: './archives.component.html',
  styleUrls: ['./archives.component.scss'],
})
export class PostArchivesEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * Type of archive whether category or tag
   */
  @Input()
  type: string;

  /**
   * List of categories or tags to-be rendered as sidebar menu-item(s)
   */
  archives: Array<Archive>;

  /**
   * List of categories or tags are matched with the search keyword
   */
  matchedSearchArchives: Array<Archive> = [];

  /**
   * List of categories or tags are shown in the dropdown menu
   */
  dropdownItems: Array<DropdownItem> = [];

  /**
   * List of categories or tags are selected from the dropdown menu
   */
  selectedArchives: Array<Archive> = [];

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSearch,
  };

  ngOnInit(): void {
    this.getAllArchives();

    this.selectedArchives = this.type === 'categories'
      ? this.post.categories
      : this.post.tags;
  }

  onChange(selectedItem: DropdownItem): void {
    const i: number = this.selectedArchives
      .findIndex(({ slug }: Archive): boolean => selectedItem.value === slug);

    if (i > -1) {
      return;
    }

    this.selectedArchives = this.selectedArchives
      .concat(this.archives.find(({ slug }: Archive): boolean => selectedItem.value === slug))
      .filter((val, i, self) => self.indexOf(val) === i);

    this.updatePostArchives(this.post.slug, this.selectedArchives.map(({ slug }: Archive): string => slug));
  }

  onBlur(): void {
    setTimeout(() => { this.dropdownItems = []; }, 800);
  }

  onKeyPress(keyword: string): void {
    this.searchAllArchives(keyword);
  }

  toggleSelectedArchive(archive: Archive): void {
    const i: number = this.selectedArchives
      .findIndex(({ slug }: Archive): boolean => archive.slug === slug);

    this.selectedArchives.splice(i, 1);

    this.updatePostArchives(this.post.slug, this.selectedArchives.map(({ slug }: Archive): string => slug));
  }

  private getAllArchives(): void {
    this.apollo.query({
      query: gql`
        {
          archives: ${this.type} { name slug }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ archives: Array<Archive> }>): Array<Archive> => result.data.archives),
    ).subscribe((archives: Array<Archive>): void => {
      this.archives = archives;
    });
  }

  private searchAllArchives(keyword: string): void {
    if (this.archives.length === 0) {
      return;
    }

    const re: RegExp = new RegExp(keyword, 'i');

    this.matchedSearchArchives =
      this.archives.reduce((result: Array<Archive>, archive: Archive): Array<Archive> => {
        if (re.test(archive.name)) {
          return result.concat(archive);
        }

        return result;
      }, []);

    this.dropdownItems =
      this.matchedSearchArchives.map(({ name, slug }: Archive): DropdownItem => ({ label: name, value: slug }));
  }

  private updatePostArchives(slug: string, archiveSlugs: Array<string>): void {
    const name: string = this.type === 'categories'
      ? 'updatePostCategories(slug: $slug, categorySlugs: $archiveSlugs)'
      : 'updatePostTags(slug: $slug, tagSlugs: $archiveSlugs)';

    this.mutate(
      `
        mutation {
          updatePostArchives: ${name} {
            ...EditablePost
          }
        }
      `,
      {
        slug,
        archiveSlugs,
      }
    ).subscribe((post: Post): void => this.changeSuccess.emit(post));
  }
}
